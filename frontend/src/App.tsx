import React from 'react';
import Transaction from "./Hooks/Transaction";
import DisplayTransaction from "./Hooks/DisplayTransactions";

function App() {
    let transactions: Transaction[] = [
        {
            Total: 1,
            Description: "test",
            TimeOfTransaction: new Date(),
            TimeOfRecord: new Date(),
            Recorder: "test",
            GroupID: 1
        }
    ];

    const getTransactions = () =>{
        fetch("/api/v1/transactions")
    }
    return (
        <div className="App">
            <div className="transactions">
                {transactions.map((transaction) => (
                    <DisplayTransaction {...transaction} />
                ))}
            </div>
            <div className="add-transaction">
                <div className="data">
                    <input defaultValue="Description"/>
                    <input defaultValue="Total"/>
                </div>
                <button>Add transaction</button>
                <button onClick={getTransactions}>Link!</button>
            </div>
        </div>
    );
}

export default App;
