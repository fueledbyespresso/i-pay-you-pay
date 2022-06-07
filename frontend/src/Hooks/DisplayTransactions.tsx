import React from 'react';
import Transaction from "./Transaction";
import '../styles/transactions.scss'

const DisplayTransaction: React.FC<Transaction> = ({Total, Description}: Transaction)=>{
    let className = 'transaction-item';

    return (
        <div className={className}>
            <h2>{Description}</h2>
            Total: ${Total}
        </div>
    )
}

export default DisplayTransaction;
