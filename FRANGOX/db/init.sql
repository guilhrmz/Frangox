CREATE DATABASE IF NOT EXISTS frangox;
USE frangox;

CREATE TABLE IF NOT EXISTS pedidos (
  id INT AUTO_INCREMENT PRIMARY KEY,
  cliente VARCHAR(100),
  produto VARCHAR(100),
  entrega VARCHAR(50),
  preco DECIMAL(10,2),
  status VARCHAR(20)
);

-- Exemplo de dados
INSERT INTO pedidos (cliente, produto, entrega, preco, status) VALUES
('ROBSON', '2X FRANGO', 'RETIRADA', 90.00, 'VENDIDO'),
('ROBSON', '3X FRANGO', 'RETIRADA', 135.00, 'PENDENTE'),
('ROBSON', '1X FRANGO', 'RETIRADA', 45.00, 'VENDIDO'),
('ROBSON', '2X FRANGO', 'RETIRADA', 90.00, 'PENDENTE');
