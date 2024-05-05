function showContactWindow() {
  const contactBtn = document.querySelector("#contact");
  const contactForm = document.querySelector("#contactForm");
  contactBtn.addEventListener("click", () => {
    contactForm.style.right = "100%";
  });
}

function addPlace(event) {
  const datesField = document.querySelector("#places");
  datesField.innerHTML += `<fieldset>
    <div class="mb-4">
        <label for="place" class="block mb-1">Location:</label>
        <input type="text" name="place" required
            class="w-full px-4 py-2 border rounded-md focus:outline-none focus:border-blue-500">
    </div>
    <div class="mb-4">
        <label for="type_of_place" class="block mb-1">Type of Date:</label>
        <select name="type_of_place" required
            class="w-full px-4 py-2 border rounded-md focus:outline-none focus:border-blue-500">
            <option value="">Select Type</option>
            <option value="Meal">Meal</option>
            <option value="Drink">Drink</option>
            <option value="Casual">Casual</option>
            <option value="Formal">Formal</option>
            <option value="Adventure">Adventure</option>
        </select>
    </div>
    <div class="mb-4">
        <label for="cost" class="block mb-1">Cost:</label>
        <input type="text" name="cost" required
            class="w-full px-4 py-2 border rounded-md focus:outline-none focus:border-blue-500">
    </div>
    <div class="mb-4">
        <label for="split" class="block mb-1">Split (Yes/No):</label>
        <select name="split" required
            class="w-full px-4 py-2 border rounded-md focus:outline-none focus:border-blue-500">
            <option value="">Select Option</option>
            <option value="True">Yes</option>
            <option value="False">No</option>
        </select>
    </div>
</fieldset>`;
  event.preventDefault();
}

document.addEventListener("DOMContentLoaded", function () {
  const tokenInput = document.getElementById("ethnicity-input");
  const tokenList = document.getElementById("ethnicity-list");

  tokenInput.addEventListener("keydown", function (event) {
    if (event.key === "Enter" || event.keyCode === 13) {
      event.preventDefault();
      const tokenValue = tokenInput.value.trim();
      if (tokenValue) {
        const token = document.createElement("div");
        token.classList.add("token");
        token.textContent = tokenValue;
        tokenList.appendChild(token);
        tokenInput.value = "";
        token.addEventListener("click", function () {
          tokenList.removeChild(token);
        });
      }
    }
  });
});
