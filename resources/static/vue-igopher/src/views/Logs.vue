<template>
  <LogsPanel />
</template>

<script lang="ts">
import { Options, Vue } from "vue-class-component";
import { inject } from "vue";
import LogsPanel from "@/components/LogsPanel.vue";
import * as config from "@/config";
import { Astor } from "@/plugins/astilectron";

class Pager {
  items: any[];
  pagedItems: any[];
  currentPage: number;
  itemsPerPage: number;

  constructor(items: any[], itemsPerPage?: number) {
    this.itemsPerPage = itemsPerPage;
    this.init();
    this.setItems(items);
  }

  init() {
    this.pagedItems = [];
    this.currentPage = 0;
    if (this.itemsPerPage === undefined) {
      this.itemsPerPage = 25;
    }
  }

  refresh() {
    this.currentPage = 0;
    this.pagedItems = [];
    for (var i = 0; i < this.items.length; i++) {
      if (i % this.itemsPerPage === 0) {
        this.pagedItems[Math.floor(i / this.itemsPerPage)] = [this.items[i]];
      } else {
        this.pagedItems[Math.floor(i / this.itemsPerPage)].push(this.items[i]);
      }
    }
  }

  prevPage() {
    if (this.currentPage > 0) {
      this.currentPage--;
    }
  }

  nextPage() {
    if (this.currentPage < this.pagedItems.length - 1) {
      this.currentPage++;
    }
  }

  setPagination(nbElem: number) {
    this.itemsPerPage = +nbElem;
    this.refresh();
  }

  setItems(items: any[]) {
    this.items = items;
    this.refresh();
  }
}

@Options({
  title: "Logs",
  components: {
    LogsPanel,
  },
  mounted() {
    const astor: Astor = inject("astor");

    config.ready(() => {
      astor.onIsReady(() => {
        let pager: Pager;
        refreshLogs();

        function refreshLogs() {
          astor.trigger("getLogs", {}, function(message: any) {
            if (message.status === config.SUCCESS) {
              let items = JSON.parse(message.msg);
              if (pager === undefined) {
                pager = new Pager(items);
              } else {
                pager.setItems(items);
                config.Toast.fire({
                  icon: "success",
                  title: "Logs successfully refreshed!",
                  timer: 1500,
                });
              }
              bindList();
            } else {
              config.Toast.fire({
                icon: "error",
                title: message.msg,
              });
            }
            document
              .querySelector("#refreshLogsIcon")
              .classList.remove("fa-spin");
          });
        }

        function bindList() {
          var pgItems = pager.pagedItems[pager.currentPage];
          var new_tbody = document.createElement("tbody");
          var old_tbody = document.getElementById("logsList");
          new_tbody.id = old_tbody.id;
          for (var i = 0; i < pgItems.length; i++) {
            var tr = document.createElement("TR");
            for (var key in pgItems[i]) {
              var td = document.createElement("TD");
              td.appendChild(document.createTextNode(pgItems[i][key]));
              tr.appendChild(td);
            }
            new_tbody.appendChild(tr);
          }
          old_tbody.parentNode.replaceChild(new_tbody, old_tbody);
          document.getElementById("pageNumber").innerHTML = String(
            pager.currentPage + 1
          );
          updateTableInfo();
        }

        function prevPage() {
          pager.prevPage();
          bindList();
        }

        function nextPage() {
          pager.nextPage();
          bindList();
        }

        function setPagination(nbElem: number) {
          pager.setPagination(nbElem);
          bindList();
        }

        function updateTableInfo() {
          let firstElem = pager.currentPage * pager.itemsPerPage + 1;
          let lastElem = 0;
          if (pager.currentPage === pager.pagedItems.length - 1) {
            lastElem = pager.items.length;
          } else {
            lastElem = firstElem + pager.itemsPerPage - 1;
          }
          document.getElementById("dataTable_info").innerHTML =
            "Showing " +
            firstElem +
            " to " +
            lastElem +
            " of " +
            pager.items.length;
        }

        // Controllers
        /// Buttons
        document
          .querySelector("#refreshLogsBtn")
          .addEventListener("click", function() {
            document.querySelector("#refreshLogsIcon").classList.add("fa-spin");
            refreshLogs();
          });
        document
          .querySelector("#prevPageBtn")
          .addEventListener("click", function() {
            prevPage();
          });
        document
          .querySelector("#nextPageBtn")
          .addEventListener("click", function() {
            nextPage();
          });

        /// Select
        document
          .querySelector("#elementsPerPageSelect")
          .addEventListener("change", (event) => {
            let nb: number = +(event.target as HTMLTextAreaElement).value;
            setPagination(nb);
          });
      });
    });
  },
})
export default class Logs extends Vue {}
</script>

<style lang="scss"></style>
