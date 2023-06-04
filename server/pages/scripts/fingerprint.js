const fppromise = import('https://fpjscdn.net/v3/1pN1YgAKn7oD3s8Jvu2J').then(fingerprintjs => fingerprintjs.load());

fppromise.then(fp => {
  fp.get().then(result => {
    const fingerprintData = {
      visitorId: result.visitorId,
    };

    console.log(JSON.stringify(fingerprintData));

    fetch('fingerprint', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(fingerprintData)
    })
    .then(response => response.json())
    .then(data => console.log(data))
    .catch(error => console.error(error));
  });
});