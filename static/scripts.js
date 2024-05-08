/*
 * ATTENTION: The "eval" devtool has been used (maybe by default in mode: "development").
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
/******/ (() => { // webpackBootstrap
/******/ 	var __webpack_modules__ = ({

/***/ "./src/index.js":
/*!**********************!*\
  !*** ./src/index.js ***!
  \**********************/
/***/ (() => {

eval("// function showContactWindow() {\n//   const contactBtn = document.querySelector(\"#contact\");\n//   const contactForm = document.querySelector(\"#contactForm\");\n//   contactBtn.addEventListener(\"click\", () => {\n//     contactForm.style.right = \"100%\";\n//   });\n// }\n\n// function addPlace(event) {\n//   const datesField = document.querySelector(\"#places\");\n//   datesField.innerHTML += `<fieldset>\n//     <div class=\"mb-4\">\n//         <label for=\"place\" class=\"block mb-1\">Location:</label>\n//         <input type=\"text\" name=\"place\" required\n//             class=\"w-full px-4 py-2 border rounded-md focus:outline-none focus:border-blue-500\">\n//     </div>\n//     <div class=\"mb-4\">\n//         <label for=\"type_of_place\" class=\"block mb-1\">Type of Date:</label>\n//         <select name=\"type_of_place\" required\n//             class=\"w-full px-4 py-2 border rounded-md focus:outline-none focus:border-blue-500\">\n//             <option value=\"\">Select Type</option>\n//             <option value=\"Meal\">Meal</option>\n//             <option value=\"Drink\">Drink</option>\n//             <option value=\"Casual\">Casual</option>\n//             <option value=\"Formal\">Formal</option>\n//             <option value=\"Adventure\">Adventure</option>\n//         </select>\n//     </div>\n//     <div class=\"mb-4\">\n//         <label for=\"cost\" class=\"block mb-1\">Cost:</label>\n//         <input type=\"text\" name=\"cost\" required\n//             class=\"w-full px-4 py-2 border rounded-md focus:outline-none focus:border-blue-500\">\n//     </div>\n//     <div class=\"mb-4\">\n//         <label for=\"split\" class=\"block mb-1\">Split (Yes/No):</label>\n//         <select name=\"split\" required\n//             class=\"w-full px-4 py-2 border rounded-md focus:outline-none focus:border-blue-500\">\n//             <option value=\"\">Select Option</option>\n//             <option value=\"True\">Yes</option>\n//             <option value=\"False\">No</option>\n//         </select>\n//     </div>\n// </fieldset>`;\n//   event.preventDefault();\n// }\n\nString.prototype.capitalize = capitalize;\n\nfunction capitalize() {\n  return this.charAt(0).toUpperCase() + this.slice(1);\n}\n\ndocument.addEventListener(\"DOMContentLoaded\", function () {\n  const tokenInput = document.getElementById(\"ethnicity-input\");\n  const tokenList = document.getElementById(\"ethnicity-list\");\n  const tokenField = document.getElementById(\"ethnicities\");\n\n  if (tokenInput) {\n    tokenInput.addEventListener(\"keydown\", function (event) {\n      if (event.key === \"Enter\") {\n        event.preventDefault();\n        const tokenValue = tokenInput.value.trim().capitalize();\n        tokenField.value += tokenValue + \",\";\n        if (tokenValue) {\n          const token = document.createElement(\"div\");\n          token.classList.add(\"token\");\n          const tokenLabel = document.createElement(\"div\");\n          tokenLabel.textContent = tokenValue;\n          tokenLabel.classList.add(\"token-label\");\n          token.appendChild(tokenLabel);\n          const closeToken = document.createElement(\"div\");\n          closeToken.innerHTML = `<i class='bx bx-x-circle align-middle'></i>`;\n          token.appendChild(closeToken);\n          closeToken.classList.add(\"token-close\");\n          tokenList.appendChild(token);\n          tokenInput.value = \"\";\n          closeToken.addEventListener(\"click\", function () {\n            tokenList.removeChild(token);\n          });\n        }\n      }\n    });\n  }\n});\n\n\n//# sourceURL=webpack:///./src/index.js?");

/***/ })

/******/ 	});
/************************************************************************/
/******/ 	
/******/ 	// startup
/******/ 	// Load entry module and return exports
/******/ 	// This entry module can't be inlined because the eval devtool is used.
/******/ 	var __webpack_exports__ = {};
/******/ 	__webpack_modules__["./src/index.js"]();
/******/ 	
/******/ })()
;