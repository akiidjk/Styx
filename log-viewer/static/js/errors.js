let notificationTimeout;
function showNotification(message, type) {
  clearTimeout(notificationTimeout);
  const alert = document.getElementById(`alert-${type}`);
  alert.classList.remove('hidden', 'fade-out');
  alert.classList.add('fade-in');
  alert.querySelector('span').textContent = message;

  notificationTimeout = setTimeout(() => {
    alert.classList.remove('fade-in');
    alert.classList.add('fade-out');
  }, 3000);
}


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
