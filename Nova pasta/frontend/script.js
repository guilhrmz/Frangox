async function carregarPedidos() {
    const response = await fetch('http://localhost:8080/pedidos');
    const pedidos = await response.json();
  
    const container = document.getElementById('orders');
    container.innerHTML = '';
  
    pedidos.forEach(pedido => {
      const div = document.createElement('div');
      div.className = 'order';
      div.innerHTML = `
        <strong>PEDIDO: #${pedido.id.toString().padStart(5, '0')}</strong><br>
        CLIENTE: ${pedido.cliente_nome}<br>
        TIPO DE ENTREGA: ${pedido.entrega}<br>
        PRODUTO: ${pedido.quantidade}x ${pedido.produto_nome}<br>
        <div class="status">
          <div>PREÇO: R$${pedido.preco.toFixed(2)}</div>
          <button class="btn ver" onclick="verPedido(${pedido.id})">VER PEDIDO</button>
          <button class="btn ${pedido.status === 'VENDIDO' ? 'vendido' : 'pendente'}">${pedido.status}</button>
        </div>
      `;
      container.appendChild(div);
    });
  }
  
  function verPedido(id) {
    window.location.href = `detalhes.html?id=${id}`;
  }
  
  if (window.location.pathname.endsWith("index.html")) {
    window.onload = carregarPedidos;
  }
  