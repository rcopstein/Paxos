# Paxos

O módulo Multi-Paxos é composto por alguns Sub-Módulos:
- Os módulos de Paxos decreto único (SinglePaxos)
- O módulo de escolha de líder (Omega)
- O módulo de mensagens (P2PLink e PaxosMessages)
- O módulo controle (Multi-Paxos)

A especificação do módulo segue a especificação do livro Reliable and Secure Distributed Programming.

## Paxos de decreto único
O módulo de Paxos de decreto único implementa a especificação do livro.
O Paxos de decreto único utiliza um canal com o módulo multi-paxos como caixa de correio. Assim, cada instância de paxos decreto único não tem seu próprio módulo de mensagens.

## Omega
O sub-módulo omega implementa um detector de falhas simulado. O sub-módulo realiza a escolha de líder para cada instância de Multi-paxos.
Para casos de teste, a função de escolha de líder pode ser definida conforme a conveniência do uso do módulo.

## P2PLink e PaxosMessages
Utiliza uma pacotes UDP para transmitir as mensagens entre instâncias de módulos Multi-Paxos.
Cada mensagem trafegada contém informação sobre a instância de Paxos decreto único referida, e é entregue a instância Multi-paxos apropriada, através do seu endereço de rede.

## Multipaxos
O módulo multi-paxos controla e gerencia a criação de instâncias de paxos decreto único.
Multi-paxos pode propor valores, e caso detecte valores não propostos na sequência de instância de paxos decreto único, corrige estas ausências com valores nulos.

O módulo multi-paxos também controla o módulo de mensagens, servindo como serviço de correio para as instâncias de paxos decreto único.
