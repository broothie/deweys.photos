const photoFormSetup = (endpoint, success, failure) => {
  // File dropper
  let file;
  const fileDropper = document.getElementById('file-dropper');
  fileDropper.addEventListener('dragover', (event) => {
    event.preventDefault();
    event.dataTransfer.dropEffect = 'move';
  });

  fileDropper.addEventListener('drop', (event) => {
    event.preventDefault();
    file = event.dataTransfer.files[0];
    fileDropper.innerText = file.name;
  });

  // Form
  const submit = document.getElementById('submit');
  const formElement = document.getElementById('photo-form');
  submit.addEventListener('click', (event) => {
    event.preventDefault();
    const form = new FormData(formElement);

    if (!!file) form.append('photo', file);

    fetch(endpoint, { method: 'post', body: form })
      .then((response) => response.ok ? success(response) : failure(response));
  });
};

