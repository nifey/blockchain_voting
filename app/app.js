var express = require('express');
var bodyParser = require('body-parser');
var path = require('path');

const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');

const ccpPath = path.resolve(__dirname, '..', 'basic-network', 'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);

const walletPath = path.join(process.cwd(), 'wallet');
const wallet = new FileSystemWallet(walletPath);
console.log(`Wallet path: ${walletPath}`);

var app = express();

app.set('view engine', 'ejs');
app.set('views', path.join(__dirname, 'views'));

app.use(bodyParser.json());
app.use(bodyParser.urlencoded({extended: false}));

app.use(express.static(path.join(__dirname, 'client')));

app.get('/', function(req, res){
	var electionStatus = checkElectionStatus();	
	electionStatus.then(function(eStatus){

		var candidates = queryAllCandidates();
		candidates.then(function(candidateList){
			res.render('index',{
				eStatus : eStatus,
				candidates: candidateList
			});
		});
	});
});

app.get('/results', function(req, res){
	var electionStatus = checkElectionStatus();	
	electionStatus.then(function(eStatus){

	var candidates = queryAllCandidates();
	candidates.then(function(candidateList){
		var sortedList = [];
		for (var candidate in candidateList){
			sortedList.push(candidateList[candidate].Record);
		}
		sortedList.sort(function(a,b){
			return b.votes - a.votes;
		});
		res.render('results',{
			eStatus : eStatus,
			candidates: sortedList
		});
	});

	});
});

app.get('/candidates', function(req, res){
	var candidates = queryAllCandidates();
	candidates.then(function(candidateList){
		res.render('candidates',{
			candidates: candidateList
		});
	});
});

app.post('/vote', function(req, res){
	var vote = castVote(req.body.voterId,req.body.candidate);
	vote.then(function(message){
		res.render('vote',{
			message: message
		});
	});
});

app.post('/voter', function(req, res){
	var electionStatus = checkElectionStatus();	
	electionStatus.then(function(eStatus){

	var voter = queryVoter(req.body.voterId);
	voter.then(function(voterData){
		if(voterData=="Invalid"){

		res.render('voter',{
			eStatus : eStatus,
			valid : false
		});

		} else {
		var voted = "Not Voted";
		if(voterData.voted){
			voted = "Voted";
		}
		res.render('voter',{
			eStatus : eStatus,
			valid : true,
			voterId: req.body.voterId,
			vStatus: voted
		});
		}
	});

	});
});

app.listen(3000, function(){
	console.log('Server Started on Port 3000');
});

async function checkElectionStatus() {
    try {
	const userExists = await wallet.exists('user1');
	if (!userExists) {
		console.log('An identity for the user "user1" does not exist in the wallet');
		console.log('Run the registerUser.js application before retrying');
		return;
	}

	const gateway = new Gateway();

	await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });

        const network = await gateway.getNetwork('mychannel');

        const contract = network.getContract('vote');


        const result = await contract.evaluateTransaction('checkElectionStatus');
	var eStatus = result.toString();
	console.log(eStatus);
	return eStatus;

    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        process.exit(1);
    }
}

async function queryAllCandidates() {
    try {
        const userExists = await wallet.exists('user1');
        if (!userExists) {
            console.log('An identity for the user "user1" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }

        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });

        const network = await gateway.getNetwork('mychannel');

        const contract = network.getContract('vote');

        const result = await contract.evaluateTransaction('queryAllCandidates');
	var res = JSON.parse(result.toString());
	return res;

    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        process.exit(1);
    }
}

async function queryVoter(voterId) {
    try {
        const userExists = await wallet.exists('user1');
        if (!userExists) {
            console.log('An identity for the user "user1" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }

        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });

        const network = await gateway.getNetwork('mychannel');

        const contract = network.getContract('vote');

        const result = await contract.evaluateTransaction('queryVoter',voterId);
	if(result.toString()==""){
	    var res = "Invalid";
	    return res;
	} else {
	    var res = JSON.parse(result.toString())
	    return res;
	}

    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        process.exit(1);
    }
}

async function castVote(voterId, candidateId) {
    try {
        const userExists = await wallet.exists('user1');
        if (!userExists) {
            console.log('An identity for the user "user1" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }

        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });

        const network = await gateway.getNetwork('mychannel');

        const contract = network.getContract('vote');

        const result = await contract.submitTransaction('castVote', voterId, candidateId);
	return result.toString();

        await gateway.disconnect();

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}


