import * as React from "react";

interface IProps {
    status: string;
}

export class MediaPendingItem extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <li className={"upload-list-item row"}>
            <div className={"col-sm-4"}>
                <div className={"placeholder"}>
                    <span>v</span>
                </div>
            </div>
            <div className={"col-sm-8 text-center"}>
                <span className={"status"}>status: {this.props.status}</span>
            </div>
        </li>;
    }

}