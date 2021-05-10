ready(() => {
    document.addEventListener('astilectron-ready', function() {
        var pager = {};
        pagerInit();
        refreshLogs();

        function refreshLogs() {
            astilectron.sendMessage({ msg: "getLogs" }, function(message) {
                if (message.status === SUCCESS) {
                    pager.items = JSON.parse(message.msg);
                    pager.refresh();
                } else {
                    toastr.error(message.msg)
                }
                document.querySelector('#refreshLogsIcon').classList.remove('fa-spin');
            });
        }

        function bindList() {
            var pgItems = pager.pagedItems[pager.currentPage];
            var new_tbody = document.createElement('tbody');
            var old_tbody = document.getElementById("logsList");
            new_tbody.id = old_tbody.id;
            for (var i = 0; i < pgItems.length; i++) {
                var tr = document.createElement('TR');
                for (var key in pgItems[i]) {
                    var td = document.createElement('TD')
                    td.appendChild(document.createTextNode(pgItems[i][key]));
                    tr.appendChild(td);
                }
                new_tbody.appendChild(tr);
            }
            old_tbody.parentNode.replaceChild(new_tbody, old_tbody);
            document.getElementById("pageNumber").innerHTML = pager.currentPage + 1;
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

        function setPagination(nbElem) {
            pager.itemsPerPage = +nbElem;
            pager.refresh();
        }

        function updateTableInfo() {
            let firstElem = pager.currentPage * pager.itemsPerPage + 1;
            let lastElem = 0;
            if (pager.currentPage === pager.pagedItems.length - 1) {
                lastElem = pager.items.length;
            } else {
                lastElem = firstElem + pager.itemsPerPage - 1;
            }
            document.getElementById("dataTable_info").innerHTML = "Showing " + firstElem + " to " + lastElem + " of " + pager.items.length;
        }

        function pagerInit() {
            pager.pagedItems = [];
            pager.currentPage = 0;
            if (pager.itemsPerPage === undefined) {
                pager.itemsPerPage = 25;
            }
            pager.prevPage = function() {
                if (pager.currentPage > 0) {
                    pager.currentPage--;
                }
            };
            pager.nextPage = function() {
                if (pager.currentPage < pager.pagedItems.length - 1) {
                    pager.currentPage++;
                }
            };
            pager.refresh = function() {
                pager.currentPage = 0;
                pager.pagedItems = [];
                for (var i = 0; i < pager.items.length; i++) {
                    if (i % pager.itemsPerPage === 0) {
                        pager.pagedItems[Math.floor(i / pager.itemsPerPage)] = [pager.items[i]];
                    } else {
                        pager.pagedItems[Math.floor(i / pager.itemsPerPage)].push(pager.items[i]);
                    }
                }
                bindList();
            };
        }

        // Controllers
        /// Buttons
        document.querySelector("#refreshLogsBtn").addEventListener("click", function() {
            document.querySelector('#refreshLogsIcon').classList.add('fa-spin');
            refreshLogs();
        });
        document.querySelector("#prevPageBtn").addEventListener("click", function() {
            prevPage();
        });
        document.querySelector("#nextPageBtn").addEventListener("click", function() {
            nextPage();
        });

        /// Select
        document.querySelector("#elementsPerPageSelect").addEventListener("change", function() {
            setPagination(this.value);
        });
    });
});