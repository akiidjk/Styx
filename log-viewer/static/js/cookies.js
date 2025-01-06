function loadPreferences() {
  const fileName = Cookies.get('logFileName');
  const filterValue = Cookies.get('filterValue');
  const pageLength = Cookies.get('pageLength');
  return { fileName, filterValue, pageLength };
}

function savePreferences(fileName) {
  Cookies.set('logFileName', fileName, { expires: 7 });
}


function savePageLength(settings) {
  const api = this.api();
  const pageInfo = api.page.info();
  Cookies.set('pageLength', pageInfo.length, { expires: 7 });
}
