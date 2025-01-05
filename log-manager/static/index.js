function savePreferences(fileName) {
  Cookies.set('logFileName', fileName, { expires: 7 });
}

function showNotification(message, type = 'success') {
  const alertElement = document.querySelector(`.alert-${type}`);
  if (alertElement) {
    alertElement.textContent = message;
    alertElement.classList.remove('hidden');
    alertElement.classList.add('fade-in');

    setTimeout(() => {
      alertElement.classList.remove('fade-in');
      alertElement.classList.add('fade-out');
      setTimeout(() => {
        alertElement.classList.add('hidden');
        alertElement.classList.remove('fade-out');
      }, 500); // Durata dell'animazione di uscita
    }, 5000); // Tempo di visualizzazione
  }
}


function loadPreferences() {
  const fileName = Cookies.get('logFileName');
  const filterValue = Cookies.get('filterValue');
  const pageLength = Cookies.get('pageLength');
  return { fileName, filterValue, pageLength };
}

async function uploadFileAndFetchLogs(file, fileNameFromCookie = null) {
  try {
    let response;
    if (!file && fileNameFromCookie) {
      response = await axios.get('/logs', {
        params: { fileName: fileNameFromCookie }
      });
    } else if (file) {
      const formData = new FormData();
      formData.append('file', file);
      response = await axios.post('/logs', formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
      });
    } else {
      throw new Error('No file provided or stored in cookies.');
    }

    return response.data;
  } catch (error) {
    console.error('Error uploading file or fetching logs:', error);
    showNotification(`Error: ${error.message}`, 'error');
    return [];
  }
}

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
      paging: true,
      pageLength: pageLength,
      drawCallback: function (settings) {
        const api = this.api();
        const pageInfo = api.page.info();
        Cookies.set('pageLength', pageInfo.length, { expires: 7 });
      },
      lengthMenu: [[100, 200, 300, 500, 1000, -1], [100, 200, 300, 500, "1K", "All"]],
      scrollX: true,
      scrollY: '50vh',
      search: {
        regex: true,
      }
    });
  });

  showNotification('Logs loaded successfully!', 'success');
}

$('#filter').on('keyup', function () {
  const query = $(this).val().toLowerCase();
  $('#logTable').DataTable().search(query).draw();
});

document.getElementById('uploadButton').addEventListener('click', async () => {
  const fileInput = document.getElementById('logFile');
  const file = fileInput.files[0];

  if (!file) {
    showNotification('Please select a file to upload.', 'error');
    return;
  }

  savePreferences(file.name);

  const logs = await uploadFileAndFetchLogs(file);
  populateTable(logs);
});

document.addEventListener('DOMContentLoaded', async () => {
  const { fileName, filterValue } = loadPreferences();

  if (fileName) {
    const logs = await uploadFileAndFetchLogs(null, fileName);
    populateTable(logs);
  }

  if (filterValue) {
    document.getElementById('filter').value = filterValue;
  }
});

window.onerror = function (message, source, lineno, colno, error) {
  console.error(`Global Error: ${message} at ${source}:${lineno}:${colno}`);
  if (error) {
    console.error('Stack Trace:', error.stack);
  }
  showNotification(`An error occurred: ${message}`, 'error');
};

window.onunhandledrejection = function (event) {
  console.error('Unhandled Promise Rejection:', event.reason);
  showNotification(`Unhandled Promise Error: ${event.reason}`, 'error');
};
