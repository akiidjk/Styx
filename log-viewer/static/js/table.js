function populateTable(data) {
  if (data.length === 0) {
    console.error('No data to display.');
    showNotification('No data to display!', 'error');
    return;
  }

  if ($.fn.DataTable.isDataTable('#logTable')) {
    $('#logTable').DataTable().destroy();
  }

  const columns = Object.keys(data[0]).map(key => ({ title: key, data: key }));

  $(document).ready(function () {
    const savedPreferences = loadPreferences();
    const pageLength = savedPreferences.pageLength ? parseInt(savedPreferences.pageLength, 10) : 10;

    $('#logTable').DataTable({
      data: data,
      columns: columns,
      responsive: true,
      destroy: true,
      paging: true,
      pageLength: pageLength,
      drawCallback: savePageLength,
      lengthMenu: [[100, 200, 300, 500, 1000, -1], [100, 200, 300, 500, "1K", "All"]],
      scrollY: '50vh',
      search: {
        regex: true,
      }
    });

    initColumnsDropdown();
  });
  showNotification('Logs loaded successfully!', 'success');
}


$('#exportButton').on('click', function () {
  const table = $('#logTable').DataTable();
  const data = table.rows({ search: 'applied' }).data().toArray();
  const csv = Papa.unparse(data);
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = 'filtered_logs.csv';
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
});


function initColumnsDropdown() {
  table = $('#logTable').DataTable();
  const menu = $('#columnToggle');
  table.columns().every(function (index) {
    const columnTitle = table.column(index).header().innerText;
    const item = `
        <label class="label cursor-pointer">
          <span class="label-text">${columnTitle}</span>
          <input type="checkbox" data-column="${index}" checked="checked" class="checkbox" />
        </label>
      `;
    menu.append(item);

    $('#columnToggle input[type="checkbox"]').on('change', function () {
      const column = $('#logTable').DataTable().column($(this).data('column'));
      column.visible($(this).is(':checked'));
    });
  });
}
