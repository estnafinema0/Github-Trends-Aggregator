import React, { useEffect, useState } from 'react';
import { Line } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';

// Регистрируем необходимые компоненты Chart.js
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);

function App() {
  const [repos, setRepos] = useState([]);
  const [selectedRepo, setSelectedRepo] = useState(null);
  const [chartData, setChartData] = useState({});

  console.log("selectedRepo", selectedRepo);
  useEffect(() => {
    fetchTrends();

    // Устанавливаем WebSocket-соединение с сервером
    const ws = new WebSocket("ws://localhost:8080/ws");
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setRepos(data);
    };
    return () => {
      ws.close();
    };
  }, []);

  const fetchTrends = async () => {
    try {
      const res = await fetch("http://localhost:8080/trends");
      const data = await res.json();
      setRepos(data);
    } catch (err) {
      console.error("Ошибка получения данных:", err);
    }
  };

  // При выборе репозитория симулируем динамику для графика
  const handleRepoSelect = (repo) => {
    console.log("Выбран репозиторий:", repo); // Отладочный лог
    setSelectedRepo(repo);
    const labels = [];
    const dataPoints = [];
    for (let i = 0; i < 10; i++) {
      labels.push(`${i} ч.`);
      const randomStars = Math.floor(Math.random() * 100);
      dataPoints.push(repo.stars + randomStars);
      console.log(`Добавлено значение для диаграммы: ${repo.stars} + ${randomStars} = ${repo.stars + randomStars}`); // Отладочный лог
    }
    setChartData({
      labels: labels,
      datasets: [
        {
          label: 'Динамика звезд',
          data: dataPoints,
          fill: false,
          backgroundColor: 'rgb(75, 192, 192)',
          borderColor: 'rgba(75, 192, 192, 0.2)'
        }
      ]
    });
    console.log("Данные для диаграммы обновлены:", { labels, dataPoints }); // Отладочный лог
  };

  return (
    <div style={{ padding: '20px' }}>
      <h1>GitHub Trends Aggregator</h1>
      <h2>Список трендовых репозиториев</h2>
      <ul>
        {repos.map(repo => (
          <li key={repo.id} style={{ marginBottom: '10px', cursor: 'pointer' }} onClick={() => handleRepoSelect(repo)}>
            <strong>{repo.author}/{repo.name}</strong> – {repo.description} ({repo.language})
          </li>
        ))}
      </ul>
      {selectedRepo && (
        <div>
          <h3>Динамика репозитория: {selectedRepo.name}</h3>
          <Line data={chartData} />
        </div>
      )}
    </div>
  );
}

export default App;


