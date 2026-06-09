

startBtn = document.getElementById('startBtn');
const reqCountInput = document.getElementById('reqCount');
 const loader = document.getElementById('loader');
const resultsContainer = document.getElementById('resultsContainer');
startBtn.addEventListener('click', async () => { 
const count = reqCountInput.value; 
resultsContainer.innerHTML = ''; 
loader.style.display = 'block'; 
startBtn.disabled = true;

 try {
    const response = await fetch(`/api/fetch?count=${count}`);
 if (!response.ok) throw new Error('Ошибка связи с сервером Go');
    const data = await response.json(); 
    data.forEach((item, index) => { 
    const card = document.createElement('div');
    card.className = `result-card ${item.source === 'JOKE' ? 'joke' : ''}`;
    card.innerHTML = ` <div class="card-source">${item.source} API</div> <div class="card-text">${item.data}</div> `;
    card.style.animationDelay = `${index * 0.15}s`; 
    resultsContainer.appendChild(card); }); 
} catch (error) {
    alert("Произошла ошибка: " + error.message); } 
finally {loader.style.display = 'none';
     startBtn.disabled = false; 
    } 
});

