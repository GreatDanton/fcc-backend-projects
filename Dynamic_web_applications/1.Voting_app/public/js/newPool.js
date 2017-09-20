var answersContainer = document.getElementById('vote-options');
var addAnswer = document.getElementById('vote-add-option');

// add new text box on button click
addAnswer.addEventListener('click', function () {
    var optionNum = answersContainer.childElementCount + 1;
    var input = "<input type='text' placeholder='Option " + optionNum + "' />";
    answersContainer.innerHTML += input;
});


//clear error messages
var errorLabels = document.getElementsByClassName('error-message');
var submitVote = document.getElementsByClassName('btn-submit')[0];

// on button submit vote click clear error labels
submitVote.addEventListener('click', function () {
    for (var i = 0; i < errorLabels.length; i++) {
        errorLabels[i].value = "";
    }
})