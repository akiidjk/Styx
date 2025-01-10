async function uploadFileAndFetchLogs(file, fileName = null) {
  try {
    let response;
    if (!file && fileName) {
      response = await axios.get("/logs", {
        params: { fileName: fileName },
      });
    } else if (file) {
      const formData = new FormData();
      formData.append("file", file);
      response = await axios.post("/logs", formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });
    } else {
      throw new Error("No file provided or stored in cookies.");
    }

    if (fileName != null) {
      document.getElementById("title").textContent =
        `Structured Log Viewer | Logs for ${fileName}`;
    }

    const contentType = response.headers["content-type"];

    if (contentType && contentType.includes("multipart/form-data")) {
      const boundary = contentType.split("boundary=")[1];
      const data = await parseMultipartResponse(response.data, boundary);
      return {
        rows: JSON.parse(data.file),
        columns: JSON.parse(data.columns),
      };
    } else {
      return JSON.parse(response.data);
    }
  } catch (error) {
    console.error("Error response:", error.response);
    showNotification(
      `Error: ${error.response?.data?.message || error.message}`,
      "error",
    );
    return null;
  }
}

async function parseMultipartResponse(data, boundary) {
  const text = await data;
  const parts = text.split(`--${boundary}`);

  const result = {};
  parts.forEach((part) => {
    if (!part.trim() || part.includes("--\r\n")) return;

    const [headers, ...bodyParts] = part.trim().split("\r\n\r\n");
    const headerMatch = headers.match(/name="([^"]+)"/);
    if (headerMatch) {
      const fieldName = headerMatch[1];
      const body = bodyParts.join("\r\n\r\n").trim();
      result[fieldName] = body;
    }
  });

  return result;
}

function initializeDropdown() {
  const dropdownItems = document.querySelectorAll(".dropdown-item");

  if (dropdownItems.length > 0) {
    dropdownItems.forEach((item) => {
      item.addEventListener("click", (event) => {
        event.preventDefault();
        const fileName = event.target.getAttribute("data-file");
        handleFileSelection(fileName);
      });
    });
  }
}

async function handleFileSelection(fileName) {
  if (!fileName) return;

  try {
    const response = await axios.get("/logs", { params: { fileName } });
    if (response.data) {
      saveFilename(fileName);
      location.reload();
      showNotification(`Logs for ${fileName} loaded successfully.`, "success");
    }
  } catch (error) {
    console.error(`Error fetching logs for ${fileName}:`, error);
    showNotification(`Error loading logs for ${fileName}.`, "error");
  }
}

document.addEventListener("DOMContentLoaded", async () => {
  const { fileName, filterValue } = loadPreferences();

  if (fileName) {
    const logs = await uploadFileAndFetchLogs(null, fileName);
    populateTable(logs.rows, logs.columns);
  }

  if (filterValue) {
    gridApi.setFilterModel(JSON.parse(filterValue));
  }

  initializeDropdown();
});

document.getElementById("uploadButton").addEventListener("click", async () => {
  const fileInput = document.getElementById("logFile");
  const file = fileInput.files[0];

  if (!file) {
    showNotification("Please select a file to upload.", "error");
    return;
  }

  saveFilename(file.name);

  const logs = await uploadFileAndFetchLogs(file);
  location.reload();
  populateTable(logs);
});
