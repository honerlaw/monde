import * as React from "react";
import {ICommentResponse} from "../media-view-page";
import {CommentForm} from "./comment-form";
import {CheckboxButton} from "../../../common/checkbox-button";

interface IProps {
    comment: ICommentResponse;
}

export class CommentItem extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <div className={"comment"}>
            <div className={"comment-view"}>
                <div className={"comment-view-header"}>
                    <span className={"username"}>{this.props.comment.user_id}</span>
                    <span className={"created"}>{this.props.comment.created_at}</span>
                </div>
                <div className={"comment-view-text"}>
                    {this.props.comment.comment}
                </div>
                <div className={"comment-view-footer"}>
                    <CheckboxButton id={this.props.comment.id} label={"reply"}>
                        <CommentForm mediaId={this.props.comment.media_id} parentCommentId={this.props.comment.id} hidden={false}/>
                    </CheckboxButton>
                </div>
            </div>
            {this.renderChildren()}
        </div>;
    }

    private renderChildren(): JSX.Element | null {
        if (!this.props.comment.children || this.props.comment.children.length === 0) {
            return null;
        }
        return <div className={"comment-children"}>
            {this.props.comment.children.map((comment: ICommentResponse): JSX.Element => {
                return <CommentItem comment={comment}/>;
            })}
        </div>;
    }

}