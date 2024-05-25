import ChartLib, { ChartConfiguration, ChartItem, Chart, BarController, PieController, CategoryScale, LinearScale, BarElement, LineController, PointElement, LineElement, ArcElement, Colors, Title, Tooltip } from "chart.js";
import ChartDataLabels from 'chartjs-plugin-datalabels';

Chart.register(ChartDataLabels, BarController, PieController, CategoryScale, LinearScale, BarElement, LineController, PointElement, LineElement, ArcElement, Colors, Title, Tooltip);
Chart.defaults.set('plugins.datalabels', {
  anchor: 'end',
  align: 'top'
})

const tooltipDollar = (tooltipItem: any) => {
  const item = tooltipItem.raw;
  if (item === "0") return item;
  return "$" + tooltipItem.raw;
}

const MONTHS: Month[] = [
  "January",
  "February",
  "March",
  "April",
  "May",
  "June",
  "July",
  "August",
  "September",
  "October",
  "November",
  "December",
];

const getPrevMonths = (
  numberOfMonths: number = 12,
  startingDate: Date = new Date()
): [Month, Year][] => {
  let out: [Month, Year][] = [];
  let month: number = startingDate.getMonth();
  let year: number = startingDate.getFullYear();
  for (let i = 0; i < numberOfMonths; i++) {
    out.push([MONTHS[month], year]);
    month--;
    if (month < 0) {
      month = 11;
      year--;
    }
  }
  return out.reverse();
};

const renderChart = async function (
  el: HTMLElement,
  cfg: ChartConfiguration
): Promise<Chart<any>> {
  return new Chart(el as ChartItem, cfg);
};

String.prototype.capitalize = function (): string {
  return this.charAt(0).toUpperCase() + this.slice(1);
};

const getDates = async (user: string = ''): Promise<Object[]> => {
  const dates = await fetch(`/api/dates/${user}`);
  const jsonDates = await dates.json();
  return jsonDates;
};

const aggDatesByMonth = (
  dates: DateData[],
  months: [Month, Year][]
): number[] => {
  const buckets: number[] = Array(months.length).fill(0);
  for (let i = 0; i < months.length; i++) {
    const [month, year] = months[i];
    buckets[i] = dates.filter((dateData) => {
      const occurrence: Date = new Date(dateData.date);
      return (
        occurrence.getFullYear() == year &&
        occurrence.getMonth() == MONTHS.indexOf(month)
      );
    }).length;
  }
  return buckets;
};

const aggCostByMonth = (
  dates: DateData[],
  months: [Month, Year][]
): number[] => {
  const buckets: number[] = Array(months.length).fill(0);
  for (let i = 0; i < months.length; i++) {
    const [month, year] = months[i];
    buckets[i] = dates
      .filter((dateData) => {
        const occurrence: Date = new Date(dateData.date);
        return (
          occurrence.getFullYear() == year &&
          occurrence.getMonth() == MONTHS.indexOf(month)
        );
      })
      .reduce((acc, currentValue) => {
        return acc + (currentValue.cost || 0);
      }, 0);
  }
  return buckets;
};

const avgCostPerDateByMonth = (
  dates: DateData[],
  months: [Month, Year][]
): number[] => {
  const buckets: number[] = Array(months.length).fill(0);
  for (let i = 0; i < months.length; i++) {
    const [month, year] = months[i];
    const datesInMonth = dates.filter((dateData) => {
      const occurrence: Date = new Date(dateData.date);
      return (
        occurrence.getFullYear() == year &&
        occurrence.getMonth() == MONTHS.indexOf(month)
      );
    });
    const costInMonth = datesInMonth.reduce((acc, currentValue) => {
      return acc + (currentValue.cost || 0);
    }, 0);
    if (datesInMonth.length == 0) continue;
    buckets[i] = parseFloat((costInMonth / datesInMonth.length).toFixed(2));
  }
  return buckets;
};

const formatDollar = (value: number) => {
  if (value === 0) return value;
  return "$" + value.toFixed(2);
};

const getSplitData = (dates: DateData[]): number[] => {
  let numberSplit = dates.filter(date => date.split).length;
  let numberNoSplit = dates.filter(date => !date.split).length;
  // Label this [Split, NoSplit]
  return [numberSplit, numberNoSplit];
};

function addVenue(event: Event) {
  const datesField:HTMLElement | null = document.querySelector("#places");
  if (datesField) {
    datesField.innerHTML +=
    `<fieldset>
    <div class="vertical">
        <label for="place" class="block">Location:</label>
        <input type="text" name="place" id="place" required>
    </div>
    <div class="vertical">
        <label for="type_of_place" class="block">Type of Date:</label>
        <select name="type_of_place" id="type_of_place" required>
            <option value="">Select Type</option>
            <option value="Restaurant">Restaurant</option>
            <option value="Drinks">Drinks</option>
            <option value="Dessert">Dessert</option>
            <option value="Casual">Casual</option>
            <option value="Formal">Formal</option>
            <option value="Adventure">Adventure</option>
        </select>
    </div>
    <div class="vertical">
        <label for="cost" class="block">Cost:</label>
        <input type="text" name="cost" id="cost" required>
    </div>
    </fieldset>`
  }
  event.preventDefault();
}

window._dt = {
  renderChart,
  getDates,
  getPrevMonths,
  aggDatesByMonth,
  aggCostByMonth,
  avgCostPerDateByMonth,
  formatDollar,
  tooltipDollar,
  getSplitData,
  addVenue,
};
