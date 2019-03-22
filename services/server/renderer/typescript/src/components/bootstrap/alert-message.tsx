import * as React from "react";

interface IProps {
    type: "primary" | "secondary" | "danger" | "warning" | "info";
    message: string | null | undefined;
}

export class AlertMessage extends React.Component<IProps, {}> {

    public render(): JSX.Element | null {
        if (!this.props.message) {
            return null;
        }
        return <div className={`alert alert-${this.props.type}`} role="alert">
            {this.props.message}
        </div>;
    }

}
