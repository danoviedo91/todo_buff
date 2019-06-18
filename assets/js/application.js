require("expose-loader?$!expose-loader?jQuery!jquery");
require("bootstrap/dist/js/bootstrap.bundle.js");

// MODAL WINDOWS SCRIPT

// This script gets the "data-delete" attribute from the anchor tag with the trash icons.
// Then transfers this string to the "href" of the anchor element inside the modal.

// Grab the trash buttons
trashBtns = document.querySelectorAll('.js-wwc-trash-btn')

// If there is any trash buttons then...
if (trashBtns.length != 0) {

    for (let i = 0; i < trashBtns.length; i++) {
        trashBtns[i].addEventListener("click", function(event) {

            let deleteUrl = trashBtns[i].dataset.delete; // Gets the data-delete string
            let deleteBtn = document.querySelector('.js-wwc-confirm-delete-btn');
            deleteBtn.setAttribute('href', deleteUrl); // Sets the data-delete string to the href
            
        });
    }
}
