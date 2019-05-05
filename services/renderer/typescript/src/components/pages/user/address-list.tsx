import * as React from "react";
import {CheckboxButton} from "../../common/checkbox-button";
import {AddressForm} from "./address-form";
import {AddressItem} from "./address-item";
import {IAddress} from "../user-page";
import "./address-list.scss";

interface IProps {
    addresses: IAddress[];
}

export class AddressList extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <div id={"address-list"}>
            <h3>
                addresses
                <CheckboxButton id={"address-create"} label={"add"}>
                    <div className={"shadow-box"}>
                        <AddressForm address={null}/>
                    </div>
                </CheckboxButton>
            </h3>
            {this.props.addresses.map((address: IAddress): JSX.Element => {
                return <AddressItem key={address.id} address={address}/>
            })}
        </div>;
    }

}