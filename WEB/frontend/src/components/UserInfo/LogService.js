let externalErrorHandler = () => {};
let externalInformationHandler = () => {};

export function setErrorHandler(handler) {
  externalErrorHandler = handler;
}

export function showError(msg) {
  externalErrorHandler(msg);
}

export function setInformationHandler(handler) {
  externalInformationHandler = handler;
}

export function showInformation(msg) {
  externalInformationHandler(msg);
}