# Multi-Paxos, a Copsteinic Approach

O m�dulo Multi-Paxos � composto por alguns Sub-M�dulos:
- Os m�dulos de Paxos decreto �nico (SinglePaxos)
- O m�dulo de escolha de l�der (Omega)
- O m�dulo de mensagens (P2PLink e PaxosMessages)
- O m�dulo controle (Multi-Paxos)

A especifica��o do m�dulo segue a especifica��o do livro Reliable and Secure Distributed Programming.

## Paxos de decreto �nico
O m�dulo de Paxos de decreto �nico implementa a especifica��o do livro.
O Paxos de decreto �nico utiliza um canal com o m�dulo multi-paxos como caixa de correio. Assim, cada inst�ncia de paxos decreto �nico n�o tem seu pr�prio m�dulo de mensagens.

## Omega
O sub-m�dulo omega implementa um detector de falhas simulado. O sub-m�dulo realiza a escolha de l�der para cada inst�ncia de Multi-paxos.
Para casos de teste, a fun��o de escolha de l�der pode ser definida conforme a conveni�ncia do uso do m�dulo.

## P2PLink e PaxosMessages
Utiliza uma pacotes UDP para transmitir as mensagens entre inst�ncias de m�dulos Multi-Paxos.
Cada mensagem trafegada cont�m informa��o sobre a inst�ncia de Paxos decreto �nico referida, e � entregue a inst�ncia Multi-paxos apropriada, atrav�s do seu endere�o de rede.

## Multipaxos
O m�dulo multi-paxos controla e gerencia a cria��o de inst�ncias de paxos decreto �nico.
Multi-paxos pode propor valores, e caso detecte valores n�o propostos na sequ�ncia de inst�ncia de paxos decreto �nico, corrige estas aus�ncias com valores nulos.

O m�dulo multi-paxos tamb�m controla o m�dulo de mensagens, servindo como servi�o de correio para as inst�ncias de paxos decreto �nico.
