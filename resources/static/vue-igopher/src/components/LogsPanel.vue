<template>
  <div class="container-fluid">
    <div class="card shadow">
      <div class="card-header py-3">
        <p class="text-primary m-0 fw-bold">Logs</p>
      </div>
      <div class="card-body">
        <div class="row">
          <div class="col-md-6 text-nowrap">
            <div
              id="dataTable_length"
              class="dataTables_length"
              aria-controls="dataTable"
            >
              <label class="form-label d-inline-flex align-items-xxl-center"
                >Show&nbsp;<select
                  id="elementsPerPageSelect"
                  class="form-select form-select-sm custom-select custom-select-sm"
                >
                  <option value="10" selected="">10</option>
                  <option value="25" selected="">25</option>
                  <option value="50">50</option>
                  <option value="100">100</option> </select
                >&nbsp;</label
              >
            </div>
          </div>
          <div class="col-md-6">
            <div
              class="text-end text-md-end dataTables_filter"
              id="dataTable_filter"
            >
              <a
                class="btn btn-primary justify-content-end btn-icon-split"
                role="button"
                id="refreshLogsBtn"
                ><span class="text-white-50 icon"
                  ><i class="fas fa-sync-alt" id="refreshLogsIcon"></i></span
                ><span class="text-white text" id="refreshLogsSpan"
                  >Refresh</span
                ></a
              >
            </div>
          </div>
        </div>
        <div
          class="table-responsive table mt-2"
          id="dataTable"
          role="grid"
          aria-describedby="dataTable_info"
        >
          <table class="table my-0" id="dataTable">
            <thead>
              <tr>
                <th>Level</th>
                <th>Message</th>
                <th>Date</th>
              </tr>
            </thead>
            <tbody id="logsList"></tbody>
            <tfoot>
              <tr>
                <td><strong>Level</strong></td>
                <td><strong>Message</strong></td>
                <td><strong>Date</strong></td>
              </tr>
            </tfoot>
          </table>
        </div>
        <div class="row">
          <div class="col-md-6 align-self-center">
            <p
              id="dataTable_info"
              class="dataTables_info"
              role="status"
              aria-live="polite"
            >
              Showing 1 to 10 of 27
            </p>
          </div>
          <div class="col-md-6">
            <nav
              class="d-lg-flex justify-content-end dataTables_paginate paging_simple_numbers"
            >
              <ul class="pagination">
                <li class="page-item">
                  <button id="prevPageBtn" class="page-link" aria-label="Prev">
                    «
                  </button>
                </li>
                <li class="page-item active">
                  <span id="pageNumber" class="page-link">1</span>
                </li>
                <li class="page-item">
                  <button id="nextPageBtn" class="page-link" aria-label="Next">
                    »
                  </button>
                </li>
              </ul>
            </nav>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { Vue } from "vue-class-component";

export default class LogsPanel extends Vue {}
</script>

<style lang="scss"></style>
