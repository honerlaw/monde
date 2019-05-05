import * as React from "react";
import "./user-page.scss";
import {AddressList} from "./user/address-list";

export interface IAddress {
    id: string;
    type: string;
    line_one: string;
    line_two: string;
    city: string;
    state: string;
    zip_code: string;
    country: string;
}

export interface IContact {
    id: string;
    contact: string;
    type: string;
    verified: boolean;
}

interface IProps {
    contacts: IContact[];
    addresses: IAddress[];
}

export class UserPage extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        console.log(this.props.contacts);
        return <div id={"user-page"}>
            <AddressList addresses={this.props.addresses}/>
        </div>;
    }

}