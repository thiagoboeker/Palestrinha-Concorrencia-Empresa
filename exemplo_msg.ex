defmodule Account do

  # Inicia o processo da Conta
  def init([name: _name, saldo: _saldo] = opts) do
    operate(opts)
  end

  def operate([name: name, saldo: saldo] = opts) do
    # Recebe mensagens de processos externos
    receive do
      # Operacao de Saque
      {:saque, valor, sender} ->
        case saldo >= valor do
          # Saque valido
          true ->
            opts = [name: name, saldo: saldo - valor]
            send(sender, {:ok, opts, "Saque", valor})
            operate(opts)
          # Saque invalido
          false ->
            send(sender, {:invalid, opts, "Saque Invalido", valor})
            operate(opts)
        end
      # Deposito
      {:deposito, valor, sender} ->
        opts = [name: name, saldo: saldo + valor]
        send(sender, {:ok, opts, "Deposito", valor})
        operate(opts)
    end
  end

  def printop(op, name, valor) do
    IO.puts "--------------------------------------------"
    IO.puts "#{op} --> Conta: #{name} --> Valor: #{valor}"
    IO.puts "--------------------------------------------"
  end
end

defmodule Main do

  def printsaldo(op, [name: name, saldo: saldo], valor) do
    IO.puts "--------------------------------"
    IO.puts "#{op} --> Conta: #{name} --> Valor: #{valor} --> Saldo: #{saldo}"
    IO.puts "--------------------------------"
  end

  def run() do
    # Cria um processo para a conta
    account = spawn(fn -> Account.init([name: "Thiago", saldo: 500]) end)

    # Numero de conjunto de transacoes saque/deposito
    transactions = 10000

    # Realiza as transacoes
    for _ <- 1..transactions do
      send(account, {:deposito, 500, self()})
      send(account, {:saque, 500, self()})
    end

    # Recebe as resposta no processo principal
    for _ <- 1..(2 * transactions) do
      receive do
        {:ok, opts, op, valor} -> printsaldo(op, opts, valor)
        {:invalid, opts, op, valor} -> printsaldo(op, opts, valor)
      end
    end
  end
end
