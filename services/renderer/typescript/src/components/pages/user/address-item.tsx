import * as React from "react";
import {IAddress} from "../user-page";
import {AddressForm} from "./address-form";
import {CheckboxButton} from "../../common/checkbox-button";
import "./address-item.scss";

type AddressKey = keyof IAddress;

const VISIBLE_ADDRESS_PROPERTIES: AddressKey[] = ["line_one", "line_two", "city", "state", "zip_code", "country"];

interface IProps {
    address: IAddress;
}

export class AddressItem extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <div className={"shadow-box address-item"}>
            <h4>{this.props.address.type} address</h4>
            <a href={`/address/remove/${this.props.address.id}`}>delete</a><CheckboxButton id={this.props.address.id} label={"modify"}>
                <AddressForm address={this.props.address}/>
                <div id={"address-info"}>
                    {Object.keys(this.props.address)
                        .filter((key: string) => VISIBLE_ADDRESS_PROPERTIES.indexOf(key as AddressKey) !== -1)
                        .map((key: string): JSX.Element => {
                            return <div className={"address-item-row"}>
                                <label>{key.replace("_", " ")}</label><span>{this.props.address[key as AddressKey]}</span>
                            </div>
                        })}
                </div>
            </CheckboxButton>
        </div>;
    }

}