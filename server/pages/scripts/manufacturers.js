var loading = document.getElementsByClassName("loader")[0];

loading.style.display = "none";
var btns = document.getElementsByClassName("answer");

for (var i = 0; i < btns.length; i++) {
    btns[i].addEventListener("click", changeStyleWhenClick);
}

function changeStyleWhenClick() {
    var parentElement = this.parentElement;
    var previousElement = document.querySelector('.clicked');


    if (this.classList.length <= 2) {
    this.classList.add("clicked");

  } else {
    this.classList.remove("clicked");

  }
}


var choices = [];
var clicks = [];

function changeChoice(choice, index) {
    if (clicks[index] == undefined) {
      clicks[index] = true;
    
      choices.push(choice);
    }
    else {
      clicks[index] = undefined;
      let indexOfChoicesElement = choices.indexOf(choice);
      choices.splice(indexOfChoicesElement, 1);
    }
}


function sendRequest() {
    hideAllAndShowLoading();

  
    const data = JSON.stringify({manufacturers: choices});
const url ='http://localhost:8080/selection/manufacturers';
fetch(url, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: data
})
.then(response => {
  if (response.ok) {
    window.location.href = "http://localhost:8080/selection";
  } else {
    throw new Error('Ошибка HTTP: ' + response.status);
  }
})
.catch(error => {
  error => console.error(error)
});
}

function hideAllAndShowLoading() {
    var image = document.getElementsByClassName("backgroundImage")[0];
    image.style.display = "none";
    var question = document.getElementsByClassName("question_frame question")[0];
    question.style.display = "none";
    var buts = document.querySelectorAll("button");
    
    buts.forEach(function(item, i, buts) {
      item.style.display = "none";
    });
    var submit = document.querySelector("input");
    submit.style.display = "none";

    document.body.style.background = "#444444";
    var loading = document.getElementsByClassName("loader")[0];
    loading.style.display = "block";
}


