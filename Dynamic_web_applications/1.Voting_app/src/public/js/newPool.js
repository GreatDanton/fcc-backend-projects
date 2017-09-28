var optionsContainer = document.getElementById('vote-options');
var addOption = document.getElementById('vote-add-option');

// add new text box on button click
addOption.addEventListener('click', function () {
    var optionNum = optionsContainer.childElementCount + 1;
    // create new input and append it to option container
    var input = document.createElement('input')
    input.type = "text";
    input.placeholder = "Option " + optionNum;
    input.name = "option-" + optionNum;
    optionsContainer.appendChild(input);
});


//clear error messages
var errorLabels = document.getElementsByClassName('error-message');
var submitPool = document.getElementsByClassName('btn-submit')[0];

// on button submit vote click, clear error labels
submitPool.addEventListener('click', function () {
    for (var i = 0; i < errorLabels.length; i++) {
        errorLabels[i].value = "";
    }
})