var gridOptionsInstance;
let gridApi;
let columns;
let excludedColumns = [];

function populateTable(rows, columnNames) {
  columnRecord = columnNames.map((column) => {
    return {
      headerName: column,
      field: column,
      sortable: true,
      filter: true,
      resizable: true,
    };
  });

  initColumnsDropdown(columnNames);

  gridOptionsInstance = {
    rowData: rows,
    columnDefs: columnRecord,
    rowSelection: {
      mode: "multiRow",
      enableClickSelection: true,
    },
    pagination: true,
    paginationPageSize: Number(loadPreferences().pageLength) || 100,
    paginationPageSizeSelector: [100, 200, 500, 1000],
  };
  const myGridElement = document.querySelector("#logGrid");
  gridApi = agGrid.createGrid(myGridElement, gridOptionsInstance);
  columns = gridApi.getColumnDefs();

  gridApi.addEventListener("filterChanged", function () {
    const filterModel = gridApi.getFilterModel();
    saveFilterValue(JSON.stringify(filterModel));
  });

  gridApi.addEventListener("paginationChanged", function () {
    const pageSize = gridApi.paginationGetPageSize();
    savePageLength(pageSize);
  });

  showNotification("Logs loaded successfully!", "success");
}

document.getElementById("clearFilter").addEventListener("click", function () {
  gridApi.setFilterModel(null);
  saveFilterValue(null);
});

document.getElementById("exportButton").addEventListener("click", function () {
  gridApi.exportDataAsCsv({
    onlySelected: false,
    onlySelectedAllPages: false,
    suppressQuotes: false,
    fileName: "filtered_logs.csv",
  });
});

function initColumnsDropdown(columnNames) {
  const menu = $("#columnDropdown");
  columnNames.map(function (column, index) {
    const columnTitle = column;
    const item = `
        <label class="label cursor-pointer">
          <span class="label-text">${columnTitle}</span>
          <input type="checkbox" data-column="${column}" checked="checked" class="checkbox" />
        </label>
      `;
    menu.append(item);

    $('#columnDropdown input[type="checkbox"]')
      .off("change")
      .on("change", function () {
        const colId = $(this).data("column");
        const isVisible = $(this).is(":checked");
        const column = gridApi
          .getColumnDefs()
          .filter((column) => colId == column.colId);
        if (column) {
          column.hide = !isVisible;
          if (isVisible) {
            excludedColumns = excludedColumns.filter((id) => id !== colId);
          } else {
            excludedColumns.push(colId);
          }
          const filteredColumns = columns.filter(
            (col) => !excludedColumns.includes(col.colId),
          );
          gridApi.setGridOption("columnDefs", filteredColumns);
        }
      });
  });
}
