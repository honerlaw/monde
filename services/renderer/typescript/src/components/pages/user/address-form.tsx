import * as React from "react";
import {IAddress} from "../user-page";
import {InputGroup} from "../../bootstrap/input-group";
import "./address-form.scss";

interface IProps {
    address: IAddress | null;
}

export class AddressForm extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <form className={"address-form"} method={"POST"} action={"/address"}>
            {this.renderID()}
            <InputGroup label={"type"}
                        name={"type"}
                        type={"text"}
                        placeholder={"type"}
                        value={this.getValue("type")}/>
            <InputGroup label={"line_one"}
                        name={"line_one"}
                        type={"text"}
                        placeholder={"line one"}
                        value={this.getValue("line_one")}/>
            <InputGroup label={"line_two"}
                        name={"line_two"}
                        type={"text"}
                        placeholder={"line two"}
                        value={this.getValue("line_two")}/>
            <InputGroup label={"city"}
                        name={"city"}
                        type={"text"}
                        placeholder={"city"}
                        value={this.getValue("city")}/>
            <InputGroup label={"state"}
                        name={"state"}
                        type={"text"}
                        placeholder={"state / province / region"}
                        value={this.getValue("state")}/>
            <InputGroup label={"zip_code"}
                        name={"zip_code"}
                        type={"text"}
                        placeholder={"zip code"}
                        value={this.getValue("zip_code")}/>
            <InputGroup label={"country"}
                        name={"country"}
                        type={"text"}
                        placeholder={"country"}
                        value={this.getValue("country")}/>
            <button type={"submit"} className={"btn btn-primary"}>save</button>
        </form>
    }

    private renderID(): JSX.Element | null {
        if (this.props.address) {
            return <input type={"hidden"} name={"id"} value={this.props.address.id}/>;
        }
        return null;
    }

    private getValue(key: keyof IAddress): string | "" {
        if (this.props.address) {
            return this.props.address[key];
        }
        return "";
    }

}