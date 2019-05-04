import * as React from "react";
import {IAddress} from "../user-page";
import {AddressForm} from "./address-form";
import {CheckboxButton} from "../../common/checkbox-button";

interface IProps {
    address: IAddress;
}

export class AddressItem extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <div className={"address-item"}>
            <div className={"address-item-row"}>
                <span>type</span><span>{this.props.address.type}</span>
            </div>
            <div className={"address-item-row"}>
                <span>line one</span><span>{this.props.address.line_one}</span>
            </div>
            <div className={"address-item-row"}>
                <span>line two</span><span>{this.props.address.line_two}</span>
            </div>
            <div className={"address-item-row"}>
                <span>city</span><span>{this.props.address.city}</span>
            </div>
            <div className={"address-item-row"}>
                <span>state</span><span>{this.props.address.state}</span>
            </div>
            <div className={"address-item-row"}>
                <span>zip code</span><span>{this.props.address.zip_code}</span>
            </div>
            <div className={"address-item-row"}>
                <span>country</span><span>{this.props.address.country}</span>
            </div>
            <CheckboxButton id={this.props.address.id} label={"modify"}>
                <AddressForm address={this.props.address}/>
            </CheckboxButton>
        </div>;
    }

}