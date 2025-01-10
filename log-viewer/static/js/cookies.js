function loadPreferences() {
  const fileName = Cookies.get("logFileName");
  const filterValue = Cookies.get("filterValue");
  const pageSize = Cookies.get("pageSize");
  return { fileName, filterValue, pageLength: pageSize };
}

function saveFilename(fileName) {
  Cookies.set("logFileName", fileName, { expires: 7 });
}

function saveFilterValue(filterValue) {
  Cookies.set("filterValue", filterValue, { expires: 7 });
}

function savePageLength(size) {
  Cookies.set("pageSize", size, { expires: 7 });
}
