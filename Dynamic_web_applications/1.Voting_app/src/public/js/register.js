// this script is error checking before the final post request
// that registers new user

// define password input fields and error message labels
var pass = document.getElementById('password');
var passConfirm = document.getElementById('password-confirm');
var labelPassConfirm = document.getElementById('e-password-confirm');


var timeout;
// compare passwords when typing in password confirmation input
passConfirm.addEventListener('keyup', function () {
    clearTimeout(timeout);
    // execute compare function when user stops timing for 0.5s
    timeout = setTimeout(function () {
        comparePasswords();
    }, 500);
});

// when passConfirm loses focus, check for password match
passConfirm.addEventListener('blur', function () {
    comparePasswords()
});


// comparePasswords check if passwords in password and
// password confirm match. Displays error message if they do not
function comparePasswords() {
    passValue = pass.value;
    passConfirmValue = passConfirm.value;

    if (passValue !== passConfirmValue) {
        labelPassConfirm.innerHTML = "Passwords do not match";
        passConfirm.classList.add('error')
        return;
    }
    // if passwords match clear label field
    labelPassConfirm.innerHTML = "";
    passConfirm.classList.remove('error');
}