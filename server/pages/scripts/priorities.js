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
    if (clicks[index] == false || clicks[index] == undefined) {
      clicks[index] = true;
      
      choices.push(choice);
      addPriority(index);
    }
    else {
      clicks[index] = false;
      let indexOfChoicesElement = choices.indexOf(choice);
      choices.splice(indexOfChoicesElement, 1);
      deletePriority(index);
    }
}

var priority = 0;
var priority_variables = document.getElementsByClassName("priority");

var priorities = new Map()
for (var i = 0; i < priority_variables.length; ++i) {
    priorities.set(priority_variables[i], "")
}

var j = -1;
var deleted_priorities = [];
var priority_values = [1, 2, 3, 4, 5, 6];
var left_priorities = [1, 2, 3, 4, 5, 6];

function set_initial_values() {
    deleted_priorities = [];
    left_priorities = [1, 2, 3, 4, 5, 6];
    j = -1;
}

function addPriority(index) {
    // если удален весь выбор, вернемся к началу
  if (deleted_priorities.length == 6) {
      set_initial_values();
  }
  if (deleted_priorities.length != 0) {
      priority = deleted_priorities.shift();
  } else {
      priority = left_priorities.shift();
  }

    priorities.set(priority_variables[index], priority);
    priority_variables[index].innerHTML = priority + ':';

}

// если значение приоритета уже было, найти наименьшее неиспользованное значение приоритета
function deletePriority(index) {
  // проверяем было ли использовано значение
  // и находим наименьшее неиспользованное значение 
  // может быть пользователь удалил больше одного элемента, создаем массив для хранения удаленных элементов
  j += 1;
  deleted_priorities[j] = priorities.get(priority_variables[index]);

  let unique_deleted_priorities = deleted_priorities.filter((element, index) => {
    return deleted_priorities.indexOf(element) === index;
  });

  deleted_priorities = unique_deleted_priorities.sort();

  priorities.set(priority_variables[index], "")
  priority_variables[index].innerHTML = '';
}


function sendRequest() {

  const data = { priorities: choices };
  const url = 'http://localhost:8080/selection/priorities';
  fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  })
  .then(response => {
    if (response.ok) {
      return response.json();
    } else {
      throw new Error('Ошибка HTTP: ' + response.status);
    }
  })
  .then(data => {
    console.log(data);
    window.location.href = "http://localhost:8080/selection/price";
  })
  .catch(error => console.error(error));

}

