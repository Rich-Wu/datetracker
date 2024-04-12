function showContactWindow() {
    const contactBtn = document.querySelector("#contact");
    const contactForm = document.querySelector("#contactForm");
    contactBtn.addEventListener("click", () => {
        contactForm.style.right = "100%";
    })
}

function addPlace(event) {
    const datesField = document.querySelector("#places");
    datesField.innerHTML +=
        `<fieldset>
    <div class="mb-4">
        <label for="place" class="block mb-1">Location:</label>
        <input type="text" id="place" name="place" required
            class="w-full px-4 py-2 border rounded-md focus:outline-none focus:border-blue-500">
    </div>
    <div class="mb-4">
        <label for="type_of_place" class="block mb-1">Type of Date:</label>
        <select id="type_of_place" name="type_of_place" required
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
        <input type="text" id="cost" name="cost" required
            class="w-full px-4 py-2 border rounded-md focus:outline-none focus:border-blue-500">
    </div>
    <div class="mb-4">
        <label for="split" class="block mb-1">Split (Yes/No):</label>
        <select id="split" name="split" required
            class="w-full px-4 py-2 border rounded-md focus:outline-none focus:border-blue-500">
            <option value="">Select Option</option>
            <option value="True">Yes</option>
            <option value="False">No</option>
        </select>
    </div>
</fieldset>`
    event.preventDefault();
}