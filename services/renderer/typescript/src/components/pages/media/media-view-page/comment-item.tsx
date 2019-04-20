import * as React from "react";
import {ICommentResponse} from "../media-view-page";
import {CommentForm} from "./comment-form";

interface IProps {
    comment: ICommentResponse
}

export class CommentItem extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <div className={"comment"}>
            <div className={"comment-view"}>
                {this.props.comment.comment}
            </div>
            <CommentForm mediaId={this.props.comment.media_id} parentCommentId={this.props.comment.id}/>
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