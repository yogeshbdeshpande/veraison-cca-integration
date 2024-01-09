# Goals

Our goal is for CVM hardware attestations of Google-provided TCB to be linked to
1. verifiable authenticity: signed corim measurements
2. auditable measurements: the signed measurement also points to a supply chain transparency report for the measured binary. A document or software package we publish shows how to calculate the measurement from the binary, and the transparency report binds the binary to an auditable source tree at commit X built with toolchain container Y, signed by an organizationally endorsed builder key that the build follows SLSA L3 operational security requirements.

and ephemeral claims, e.g.,
3. vulnerability reporting: short-lived certificates of firmware status, like "has the most up to date security version number" or "is subject to CVE xyz. Restart your instance to get on the latest version". This could be modeled as a CoRIM endorsement of a claim like "uptodate as of TIMESTAMP".
4. Platform security reporting: short-lived certificates of platform firmware status, like "you can be sure that an attestation's TCB is >= x anywhere in the fleet"

# Supply chain standards

There are further things you can do with the transparency report like non-repudiation through hosting the build attestation with a transparency service that has append-only logs after identity-proofing, but there seems to be a fundamental disagreement between the IETF SCITT workstream and the sigstore.dev project on how to achieve that, since SCITT wants W3C DID identities, and sigstore.dev is already built to use OIDC. I don't know how that all is supposed to be consonant with RATS, since there's nothing in the corim or eat documents about using DID for identities. There is EAT binding to OIDC tokens though. Is there anyone in the RATS group that is participating in the SCITT effort that can explain this to me?

The SLSA provenance schema itself is defined in terms of a completely different attestation format called in-toto (https://in-toto.io), and communication I've had with them is that in-toto should be considered an alternative carrier format to CoRIM to fit into the RATS framework. If we want to link the reference measurement to an in-toto attestation, that seems like something verifier-specific that we'd need to say, "hey if you want to ensure the firmware measurement is not only signed, but built transparently, then download the SLSA attestation in dependent-rims. By the way if there's more than one thing in dependent-rims, you can understand any url with prefix X to be a firmware build attestation from Google" which is an unfortunate complexity.

# Modeling CVM attestation

I'm trying to understand how to fit all these goals into the RATS framework such that we can propose extensions to open source verifiers that aren't overly burdensome or highly specific to each particular package we want to provide reference values (and provenances) for.

In terms of the firmware measurement, we can deliver a CoRIM through a UEFI variable pointed to by the NIST SP 800-155 unmeasured event, and we can give the AMD SEV-SNP VCEK certificate through extended guest request, but everything else seems to be up to the verifier to collect independently of the VM. 

## Evidence collection

The way we're collecting attestations at the moment is through a recommended software package https://github.com/google/go-tpms-tools that wraps up a vTPM quote with a TEE quote and supporting certificates as a protocol buffer. I'm not clear if this unsigned bundling process should be modeled as any particular thing in the RATS framework. I think we're working with the "passport model" of attestation.

I don't have a sense of how the WG foresees how evidence should be bundled to give to a verifier. I'm working from a vendor-specific understanding at the moment that whatever verifier service you use, you need to use their format and API, but of course ideally I'd like this to be more of a federated arena where you can have n-of-k verifiers say some evidence matches policy, and the evidence is not too vendor-specific for that to be out of the question.

## CVM Profiles

Whereas Google has an attestation verifier service that generates an EAT with its own claims bound to an OIDC token (for the Confidential Space product), we'd like to use more standard claims, like AMD SEV-SNP measurement, Intel TDX MRTD, etc. Azure's attestation service has their own x-ms-* extensions for this that will hopefully help AMD and Intel align on how claims should be proposed for the CoRIM format.

Supposing we do get profiles from Intel and AMD for their CVM attesting environments (more below), those environments sign quotes / attestation reports that serve as evidence for the claims defined in those profiles.

I as a Reference Value Provider want to be able to provide a document that says something that covers 1 and 2 up front like, "if your AMD measurement is contained in {x, ...} or your TDX measurement is contained in {y, ...}, then you're running Google-authentic virtual firmware with security version n. The firmware this measures can be found at z".

My understanding of how to do this is for the firmware CoRIM to have a single CoMID tag and the SLSA provenance linked from dependent RIMs.
The CoMID tag will have lang: en-us, tag-identity: some-uuid we generate before signing, and triples-map containing some reference triples.
We have reference triples for both AMD and TDX by using different environment-maps with different class fields.
AMD SEV-SNP's class is up to AMD to profile, but let's just say it's a class-id for the VCEK extension oid prefix 1.3.6.1.4.1.3704.1.1. The measurement-map for this can have an mkey or not. If we had one, I'm unsure if it's something that Google would define or if it's still up to AMD. If Google, we could use a uuid that stands for Google Compute Engine?
The mval as a measurement-values-map would then contain our AMD firmware svn, and AMD profile-specific claims, but I think we'd just give the measurements and some form of acceptable policy specification. We just have one guest policy we apply everywhere, but if that changes we probably need the AMD profile to have expressions like ranges, lower- and upper-bounds for policy components.
For Intel, they'd need a similar profile for the TDREPORT components as claims.

I say measurements and not measurement even though we're talkabout about a single firmware binary because both AMD and TDX can have multiple measurements based on the VM construction, such as how many vCPUs it launched with (AMD has VMSAs and Intel has TDVPS).
For now our security version number matches what we measure as EV_S_CRTM_VERSION in PCR0, but that may change if there are technology-specific changes.

As far as I understand, the Intel profile for CoRIM only supports the boot chain up to the quoting enclave (QE) in terms of its TCB version, but the profile does not describe the QE as its own attesting environment for SGX enclave or TDX VM. The attesting key is generated in the QE and is signed by the PCE's hold on the PCK, which is per-machine-per-TCB (ppid + pceid). The quote wraps around the attesting key's signature for verification against their non-x.509 format.

AMD similarly does not have a profile for the SNP firmware as an attesting environment for an SEV-SNP VM.

# Evidence Appraisal

Setting aside evidence formats, I want to really understand how we go from a signed CoRIM and a CVM attestation to an attestation result (which I'll handwave is some JWT representation of the accepted claims).

We somehow get the VCEK or PCK certificate and attestation report / quote, and the Google firmware CoRIM to the verifier. The verifier can verify the evidence back to the manufacturer with this forwarded (or cached) collateral and introduce every quote/report field as claims of the target environment.
Let's say Google's code signing root key is in the trust anchor, so any CoRIM we sign is trusted.

If I read the CoRIM document about matching reference values against evidence, the document starts talking about conditional endorsements instead, which are a different triple from reference-value-triples. We discussed a little in the Github issues that reference values are a special kind of endorsement, but it's still jarring. It goes on to say that reference-value-triples is essentially redundant with the conditional-endorsement-triples, but you can use either. Then there's "In the reference-triple-record these are encoded together. In other triples multiple Reference Values are represented more compactly by letting one environment-map apply to multiple measurement-maps." 

It seems "Conditional Endorsement" is philosophical, and "conditional-endorsement-triples" is one implementation of the idea, and "Reference Value" is philosophical, but "reference-value-triples" is one implementation of the idea. Another implementation of "Reference Value" as an mkey of a "conditional-endorsement-triples", and the mval is more explicit about what claims are introduced. For "reference-value-triples", I don't see any explicit representation of a claim, rather, reference-value-triples lead to "authorized-by" getting added to fields of an Accepted Claim Set entry which itself is only a conceptual type to help understand appraisal, but not an actual claim itself–is this where a profile-defined claim needs to clarify meaning? I see this authorized-by as conceptually different from the optional field of a measurement-map, since that is from the CoRIM that I've signed and isn't part of an attestation result representation.

If I'm looking at a JWT with an AMD profile claim about the measurement value, I'd like another claim that the measurement value is signed by Google, or a stronger claim that the measurement value was signed by a trusted source, and the build provenance is [some google URL to the SLSA provenance].
Again though, if at all possible these claims should appeal more broadly than just Google.
