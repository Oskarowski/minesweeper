/**
 * Toggles the visibility of the input field with the given id based on the value of the given checkbox.
 * If the checkbox is checked, the input field is hidden; otherwise, it is shown.
 * @param {string} fieldId - The html id of the field to toggle.
 * @param {HTMLInputElement} checkbox - The checkbox element that toggles the input field.
 */
function toggleFieldVisibility(fieldId, checkbox) {
    const field = document.getElementById(fieldId);
    if (checkbox.checked) {
        // Clear the value and hide the field when checkbox is checked
        field.value = "";
        adjustMinesInputFieldRange();
        field.style.display = "none";
    } else {
        field.style.display = "block";
    }
}

function adjustMinesInputFieldRange() {
    const gridSizeValue = document.getElementById(
        "grid-size-input-field",
    ).value;
    const minesInputField = document.getElementById("mines-input-field");

    if (gridSizeValue) {
        const maxMines = Math.floor(gridSizeValue ** 2 * 0.8);
        minesInputField.min = 1;
        minesInputField.max = maxMines;
        minesInputField.placeholder = `Enter number of mines (max: ${maxMines})`;
    } else {
        minesInputField.min = 1;
        minesInputField.max = 350;
        minesInputField.placeholder = "Enter number of mines";
    }
}

window.addEventListener("load", () => {
    const gridSizeInputField = document.getElementById("grid-size-input-field");

    const maxGridSizeBasedOnScreenSize = window.innerWidth < 768 ? 10 : 22;
    gridSizeInputField.max = maxGridSizeBasedOnScreenSize;
});
