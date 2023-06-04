// количество номеров страниц, видимых пользователю одновременно
var sample = 10;

// сокрытие всех номеров страниц, кроме первых (в данном случае 10) страниц. Переменная sample равна 10.
for (var i = sample; i < document.getElementsByClassName("number").length - 1; i++) {
    document.getElementsByClassName("number")[i].style.display = "none";
}

// сокрытие кнопки с левой стрелкой. 
document.getElementsByClassName("arrow")[0].style.display = "none";


// указатель, с помощью которого отсчитываются номера, которые нужно сокрыть или отобразить
var pointer = 0;

// функция, показывающая следующие страницы в количестве sample и скрывающая текущие sample страниц. Вызывается нажатием кнопки с правой стрелкой.
function show_next_numbers() {

    
    // сокрытие кнопки с правой стрелкой, когда показываются последние номера страниц.
    //(если показываеются последние номера страниц, то кнопка с правой стрелкой, показывающая следующие страницы, уже не нужна, так как следующих страниц нет) 
    // Первый операнд условия - если количество последних номеров страниц меньше sample
    // Второй операнд условия - если количество последних номеров страниц равно sample
    if (pointer + sample + document.getElementsByClassName("number").length % sample == document.getElementsByClassName("number").length || 
            pointer + 2*sample == document.getElementsByClassName("number").length - 1) {
        
        // сокрытие кнопки с правой стрелкой, когда показываются последние номера страниц.
        document.getElementsByClassName("arrow")[1].style.display = "none";
    }

    // влключение кнопки с левой стрелкой один раз при нажатии на кнопку с правой стрелкой
    // условие нужно, чтобы включить только один раз, а не каждый раз при вызове данной функции
    if (pointer == 0) {
        document.getElementsByClassName("arrow")[0].style.display = "block";
    }
    
    // сокрытие текущих номеров страниц в количестве sample
    for (var i = pointer; i < pointer + sample; i++) {
        document.getElementsByClassName("number")[i].style.display = "none";
    }

    pointer += sample;

    // отображение следующих номеров страниц в количестве sample
    for (var i = pointer; i < pointer + sample; i++) {

        // проверка на последний номер. Если не будет проверки, то функция попытается отобразить следующий после последнего номер, 
        // которого не существует, и обратиться к памяти, в которой не лежит значение.
        if (i == document.getElementsByClassName("number").length - 1) {
            break;
        }

        // отображение следующего номера
        document.getElementsByClassName("number")[i].style.display = "block";
    }
}


// функция, показывающая предыдущие страницы в количестве sample и скрывающая текущие sample страниц. Вызывается нажатием кнопки с левой стрелкой.
function show_previous_numbers() {

    // сокрытие кнопки с левой стрелкой.
    //(если показываются первые номера страниц, то кнопка с левой стрелкой, показывающая предыдущие страницы, уже не нужна, так как предыдущих страниц нет) 
    if (pointer == sample){
        document.getElementsByClassName("arrow")[0].style.display = "none";
    }

    // отображение кнопки с правой стрелкой, исчезнувшей при отображении последнего ряда номеров, после того 
    // как была нажата кнопка с левой стрелкой, вызвавшей отображение предпоследнего ряда номеров
    if (pointer + document.getElementsByClassName("number").length % sample == document.getElementsByClassName("number").length || 
            pointer + sample == document.getElementsByClassName("number").length - 1) {

        // начальное значение счетчика цикла for, скрывающего элементы, когда кнопку с левой стрелкой нажимают, чтобы отобразить предпоследний ряд номеров
        var i = document.getElementsByClassName("number").length - 2;

        // отображение кнопки с правой стрелкой
        document.getElementsByClassName("arrow")[1].style.display = "block";
    } else {

        // начальное значение счетчика цикла for, скрывающего элементы, когда нажимают кнопку с левой стрелкой в остальных случаях
        var i = pointer + sample - 1;
    } 

    // сокрытие текущих номеров страниц в количестве sample
    for (; i >= pointer; i--) {
        document.getElementsByClassName("number")[i].style.display = "none";
    }

    // отображение следующих номеров страниц в количестве sample
    for (var i = pointer - 1; i >= pointer - sample; i--) {
        document.getElementsByClassName("number")[i].style.display = "block";
    }

    pointer -= sample;   
}

// скрыть вопрос
function HideAnswer() {
    var answer_button = document.getElementById("answer_button");
    // console.log(answer_button)
    answer_button.style.display = "none";
}

