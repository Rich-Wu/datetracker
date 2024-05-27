document.addEventListener("DOMContentLoaded", function () {
  const browseButtons = document.querySelectorAll(".userCard");
  if (browseButtons) {
    browseButtons.forEach((button) => {
      const browse = button.querySelector(".hover-light");
      button.addEventListener("mouseenter", (event) => {
        browse?.classList.add("hovered");
      });
      button.addEventListener("mouseleave", (event) => {
        browse?.classList.remove("hovered");
      });
    });
  }

  const tokenInput: HTMLInputElement | null = document.getElementById(
    "ethnicity-input"
  ) as HTMLInputElement;
  const tokenList: HTMLElement = document.getElementById(
    "ethnicity-list"
  ) as HTMLElement;
  const tokenField: HTMLInputElement | null = document.getElementById(
    "ethnicities"
  ) as HTMLInputElement;

  function addToken(): void {
    const tokenValue = tokenInput?.value.trim().capitalize();
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
      if (tokenInput) {
        tokenInput.value = "";
      }
      closeToken.addEventListener("click", function () {
        tokenList.removeChild(token);
      });
    }
  }

  if (tokenInput) {
    tokenInput.addEventListener("blur", function (event) {
      event.preventDefault();
      addToken();
    });
    tokenInput.addEventListener("keydown", function (event) {
      if (event.key === "Enter") {
        event.preventDefault();
        addToken();
      }
    });
  }
});
