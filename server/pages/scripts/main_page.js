
var preloader = document.getElementById("preloader");
preloader.style.display = "none";
var mark;

// показывать модели машин в соответствии с маркой
var clicked = false;
$("#mark").on("change", function() {

	mark = $( "#mark option:selected" ).text();
	

	if (!clicked) {
	  clicked = true;
	  $("#model").removeAttr("disabled");
	}
	var id = $(this).find('option:selected').attr('id');;
	
	var id = $(this).find('option:selected').attr('id');
	
	$("#model option").hide(); //hide all options from slect box
	
		$("#model option[id='" + id + "']").show(); //show that option

})




// валидация ввода цены: только числа
$(function() {
	$("input[class='form-control']").on('input', function(e) {
		$(this).val($(this).val().replace(/[^0-9]/g, ''));
	});
});




function sendRequest() {
	
    const data = JSON.stringify({new: check});
	
    const url ='http://localhost:8080/main';
    var req = new XMLHttpRequest();
    req.open("POST", url, true);
    req.setRequestHeader('Content-Type', 'application/json');
    req.send(data);
}




function hideAllAndShowLoading() {
	var form = document.getElementById("usual_search")
	form.style.display = "none";
	var fuzzy_algorithm = document.getElementById("fuzzy_algorithm")
	fuzzy_algorithm.style.display = "none"
	preloader.style.display = "block";

}
