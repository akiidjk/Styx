document.getElementById('resetFilter').addEventListener('click', () => {
  $('#filter').val('');
  $('#logTable').DataTable().search('').draw();
  Cookies.remove('filterValue');
});


$('#filter').on('keyup', function () {
  const query = $(this).val().toLowerCase();
  Cookies.set('filterValue', query, { expires: 7 });
  $('#logTable').DataTable().search(query).draw();
});
