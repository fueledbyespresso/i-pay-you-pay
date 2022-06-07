export default interface Transaction {
    Total: number;
    Description: string;
    TimeOfTransaction: Date;
    TimeOfRecord: Date;
    Recorder: string;
    GroupID: number
}