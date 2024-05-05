// function showContactWindow() {
//   const contactBtn = document.querySelector("#contact");
//   const contactForm = document.querySelector("#contactForm");
//   contactBtn.addEventListener("click", () => {
//     contactForm.style.right = "100%";
//   });
// }

// function addPlace(event) {
//   const datesField = document.querySelector("#places");
//   datesField.innerHTML += `<fieldset>
//     <div class="mb-4">
//         <label for="place" class="block mb-1">Location:</label>
//         <input type="text" name="place" required
//             class="w-full px-4 py-2 border rounded-md focus:outline-none focus:border-blue-500">
//     </div>
//     <div class="mb-4">
//         <label for="type_of_place" class="block mb-1">Type of Date:</label>
//         <select name="type_of_place" required
//             class="w-full px-4 py-2 border rounded-md focus:outline-none focus:border-blue-500">
//             <option value="">Select Type</option>
//             <option value="Meal">Meal</option>
//             <option value="Drink">Drink</option>
//             <option value="Casual">Casual</option>
//             <option value="Formal">Formal</option>
//             <option value="Adventure">Adventure</option>
//         </select>
//     </div>
//     <div class="mb-4">
//         <label for="cost" class="block mb-1">Cost:</label>
//         <input type="text" name="cost" required
//             class="w-full px-4 py-2 border rounded-md focus:outline-none focus:border-blue-500">
//     </div>
//     <div class="mb-4">
//         <label for="split" class="block mb-1">Split (Yes/No):</label>
//         <select name="split" required
//             class="w-full px-4 py-2 border rounded-md focus:outline-none focus:border-blue-500">
//             <option value="">Select Option</option>
//             <option value="True">Yes</option>
//             <option value="False">No</option>
//         </select>
//     </div>
// </fieldset>`;
//   event.preventDefault();
// }

String.prototype.capitalize = capitalize;

function capitalize() {
  return this.charAt(0).toUpperCase() + this.slice(1);
}

document.addEventListener("DOMContentLoaded", function () {
  const tokenInput = document.getElementById("ethnicity-input");
  const tokenList = document.getElementById("ethnicity-list");
  const tokenField = document.getElementById("ethnicities");

  if (tokenInput) {
      tokenInput.addEventListener("keydown", function (event) {
        if (event.key === "Enter") {
          event.preventDefault();
          const tokenValue = tokenInput.value.trim().capitalize();
          tokenField.value += tokenValue + ",";
          if (tokenValue) {
            const token = document.createElement("div");
            token.classList.add("token");
            const tokenLabel = document.createElement("div");
            tokenLabel.textContent = tokenValue;
            tokenLabel.classList.add("token-label");
            token.appendChild(tokenLabel);
            const closeToken = document.createElement("div");
            closeToken.innerHTML = `<i class='bx bx-x-circle align-middle'></i>`;
            token.appendChild(closeToken);
            closeToken.classList.add("token-close");
            tokenList.appendChild(token);
            tokenInput.value = "";
            closeToken.addEventListener("click", function () {
                    tokenList.removeChild(token);
            });
          }
        }
      });
  }
});