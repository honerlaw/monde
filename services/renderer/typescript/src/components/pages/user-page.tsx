import * as React from "react";
import "./user-page.scss";
import {AddressForm} from "./user/address-form";
import {AddressItem} from "./user/address-item";
import {CheckboxButton} from "../common/checkbox-button";

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

interface IProps {
    addresses: IAddress[]
}

export class UserPage extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <div id={"user-page"}>
            <div id={"address-list"}>
                <CheckboxButton id={"address-create"} label={"add"}>
                    <AddressForm address={null}/>
                </CheckboxButton>
                {this.props.addresses.map((address: IAddress): JSX.Element => {
                    return <AddressItem address={address}/>
                })}
            </div>
        </div>;
    }

}