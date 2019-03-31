import * as React from "react";
import {IUploadForm} from "./upload-list-page/upload-form";
import {PendingUploadItem} from "./upload-list-page/pending-upload-item";
import {UploadItem} from "./upload-list-page/upload-item";
import {IGlobalProps} from "../../global-props";

export interface IMediaVideoResponse {
    type: string;
    width: number;
    height: number;
    url: string;
}

export interface IMediaResponse {
    id: string;
    title: string;
    description: string;
    transcoding_status: string;
    hashtags: string[];
    is_published: boolean;
    can_publish: boolean;
    thumbnails: string[];
    videos: IMediaVideoResponse[];
}

interface IProps extends IGlobalProps {
    uploads: IMediaResponse[];
    uploadForm: IUploadForm;
}

/**
 * @todo
 * - display thumbnail after upload so the user can see what it is
 */
export class UploadListPage extends React.Component<IProps, {}> {

    public render(): JSX.Element {
        return <div id={"upload-list-page"}>
            <ol className={"upload-list"}>
                {this.props.uploads.map((upload: IMediaResponse): JSX.Element => {
                    if (upload.transcoding_status !== "Complete") {
                        return <PendingUploadItem key={upload.id} status={upload.transcoding_status}/>;
                    }
                    return <UploadItem key={upload.id} upload={upload}/>;
                })}
            </ol>
        </div>;
    }

}
