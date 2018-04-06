<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width,
    user-scalable=no, initial-scale=1.0, maximum-scale=1.0,
minimum-scale=1.0">
    <link rel="stylesheet"
href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css">
    <link rel="stylesheet" href="styles.css">
            <title> ppgc </title>
</head>

<h1> {{.Title}} </h1>


<h4> image count: {{.ImageCount}} </h4>

<div style="margin:30px; background-color: #e2e0e6">

<h3> Upload file </h3>
<!--<form enctype="multipart/form-data" action="/upload" method="post">-->
    <!--<input type="file" name="uploadfile" />-->
    <!--<input type="hidden" name="token" value="{{.}}"/>-->
    <!--<input type="submit" value="upload" />-->
<!--</form>-->

</div>

   <form action="/upload" method="post" enctype="multipart/form-data" class="uploadForm">
        <input class="uploadForm__input" type="file" name="file" id="inputFile" accept="image/*">
        <label class="uploadForm__label" for="inputFile">
            <i class="fa fa-upload uploadForm__icon"></i> Select a file
        </label>
    </form>
    <div class="notification" id="alert"></div>
    <script src="https://unpkg.com/axios/dist/axios.min.js"></script>
    <script src="api.js"></script>
    <script src="app.js"></script>

<a href="/agg"> new agg </a>



