document.addEventListener("DOMContentLoaded", function () {
  const tokenInput: HTMLInputElement | null = document.getElementById(
    "ethnicity-input"
  ) as HTMLInputElement;
  const tokenList: HTMLElement = document.getElementById(
    "ethnicity-list"
  ) as HTMLElement;
  const tokenField: HTMLInputElement | null = document.getElementById(
    "ethnicities"
  ) as HTMLInputElement;

  if (tokenInput) {
    tokenInput.addEventListener("keydown", function (event) {
      if (event.key === "Enter") {
        event.preventDefault();
        const tokenValue = tokenInput.value.trim().capitalize();
        if (tokenField) {
          tokenField.value += tokenValue + ",";
        }
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