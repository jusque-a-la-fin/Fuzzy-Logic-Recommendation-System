function sendRequest() {
    const price = {
      minPrice: document.getElementById('lst_price').value,
      maxPrice: document.getElementById('hst_price').value,
      
    };
    console.log("ффффффффф")
    
    const data = JSON.stringify(price);
    const url ='http://localhost:8080/selection/price';
    
    fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: data
    })
    .then(response => {
      if (response.ok) {
        window.location.href = "http://localhost:8080/selection/manufacturers";
      } else {
        throw new Error('Ошибка HTTP: ' + response.status);
      }
    })
    .catch(error => console.error(error));
  }