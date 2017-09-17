var answersContainer = document.getElementById('vote-answers');
var addAnswer = document.getElementById('vote-add-answer');

// add new text box on button click
addAnswer.addEventListener('click', function () {
    var optionNum = answersContainer.childElementCount;
    var input = "<input type='text' placeholder='Option " + optionNum + "' />";
    answersContainer.innerHTML += input;
});