import Chart, { ChartConfiguration, ChartItem } from "chart.js/auto";
import ChartDataLabels from 'chartjs-plugin-datalabels';

Chart.register(ChartDataLabels);
Chart.defaults.set('plugins.datalabels', {
  anchor: 'end',
  align: 'top'
})

const tooltipDollar = (tooltipItem: any) => {
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
  return "$" + value.toFixed(2);
}

window._dt = {
  renderChart,
  getDates,
  getPrevMonths,
  aggDatesByMonth,
  aggCostByMonth,
  avgCostPerDateByMonth,
  formatDollar,
  tooltipDollar
};
